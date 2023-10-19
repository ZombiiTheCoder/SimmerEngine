import argparse
import subprocess
from typing import *
import tempfile
import copy
import os
import shutil
import sys
import re
import configparser
from types import SimpleNamespace
from textwrap import dedent

def read_cfg():
    src_dir = os.path.dirname(os.path.abspath(__file__))
    cfg = configparser.ConfigParser(allow_no_value=True)
    cfgs = cfg.read("engine-bolt-wrapper.ini")
    if not cfgs:
        cfgs = cfg.read(os.path.join(src_dir, "engine-bolt-wrapper.ini"))
    assert cfgs, f"engine-bolt-wrapper.ini is not found in {os.getcwd()}"

    def get_cfg(key):
        # if key is not present in config, assume False
        if key not in cfg["config"]:
            return False
        # if key is present, but has no value, assume True
        if not cfg["config"][key]:
            return True
        # if key has associated value, interpret the value
        return cfg["config"].getboolean(key)

    d = {
        # BOLT binary locations
        "BASE_BOLT": cfg["config"]["base_bolt"],
        "CMP_BOLT": cfg["config"]["cmp_bolt"],
        # optional
        "VERBOSE": get_cfg("verbose"),
        "KEEP_TMP": get_cfg("keep_tmp"),
        "NO_MINIMIZE": get_cfg("no_minimize"),
        "RUN_SEQUENTIALLY": get_cfg("run_sequentially"),
        "COMPARE_OUTPUT": get_cfg("compare_output"),
        "SKIP_BINARY_CMP": get_cfg("skip_binary_cmp"),
        "TIMING_FILE": cfg["config"].get("timing_file", "timing.log"),
    }
    if d["VERBOSE"]:
        print(f"Using config {os.path.abspath(cfgs[0])}")
    return SimpleNamespace(**d)

# perf2bolt mode
PERF2BOLT_MODE = ["-aggregate-only", "-ignore-build-id"]

# boltdiff mode
BOLTDIFF_MODE = ["-diff-only", "-o", "/dev/null"]

# options to suppress binary differences as much as possible
MINIMIZE_DIFFS = ["-bolt-info=0"]

# bolt output options that need to be intercepted
BOLT_OUTPUT_OPTS = {
    "-o": "BOLT output binary",
    "-w": "BOLT recorded profile",
}

# regex patterns to exclude the line from log comparison
SKIP_MATCH = [
    "BOLT-INFO: BOLT version",
    r"^Args: ",
    r"^BOLT-DEBUG:",
    r"BOLT-INFO:.*data.*output data",
    "WARNING: reading perf data directly",
]


def run_cmd(cmd, out_f, cfg):
    if cfg.VERBOSE:
        print(" ".join(cmd))
    return subprocess.Popen(cmd, stdout=out_f, stderr=subprocess.STDOUT)

def run_bolt(bolt_path, bolt_args, out_f, cfg):
    p2b = os.path.basename(sys.argv[0]) == "perf2bolt"  # perf2bolt mode
    bd = os.path.base(sys.argv[0]) == "siverengine-boltdiff"
    hm = sys.argv[1] == "heatmap"  # heatmap mode
    cmd = ["/usr/bin/time", "-f", "%e %M", bolt_path] + bolt_args
    if p2b:
        # -ignore-build-id can occur at most once, hence remove it from cmd
        if "-ignore-build-id" in cmd:
            cmd.remove("-ignore-build-id")
        cmd += PERF2BOLT_MODE
    elif bd:
        cmd += BOLTDIFF_MODE
    elif not cfg.NO_MINIMIZE and not hm:
        cmd += MINIMIZE_DIFFS
    return run_cmd(cmd, out_f, cfg)

def prepend_dash(args: Mapping[AnyStr, AnyStr]) -> Sequence[AnyStr]:
    """
    Accepts parsed arguments and returns flat list with dash prepended to
    the option.
    Example: Namespace(o='test.tmp') -> ['-o', 'test.tmp']
    """
    dashed = [("-" + key, value) for (key, value) in args.items()]
    flattened = list(sum(dashed, ()))
    return flattened


def replace_cmp_path(tmp: AnyStr, args: Mapping[AnyStr, AnyStr]) -> Sequence[AnyStr]:
    """
    Keeps file names, but replaces the path to a temp folder.
    Example: Namespace(o='abc/test.tmp') -> Namespace(o='/tmp/tmpf9un/test.tmp')
    Except preserve /dev/null.
    """
    replace_path = (
        lambda x: os.path.join(tmp, os.path.basename(x))
        if x != "/dev/null"
        else "/dev/null"
    )
    new_args = {key: replace_path(value) for key, value in args.items()}
    return prepend_dash(new_args)

def preprocess_args(args: argparse.Namespace) -> Mapping[AnyStr, AnyStr]:
    """
    Drop options that weren't parsed (e.g. -w), convert to a dict
    """
    return {key: value for key, value in vars(args).items() if value}

def write_to(txt, filename, mode="w"):
    with open(filename, mode) as f:
        f.write(txt)
    
def wait(proc, fdesc):
    proc.wait()
    fdesc.close()
    return open(fdesc.name)

