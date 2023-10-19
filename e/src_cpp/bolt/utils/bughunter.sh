BOLT=${BOLT:=SimerEngine-bolt}

ulimit -c 0
set -o pipefail

if [[ -n "$VERBOSE" ]]; then
    set -x
fi

if [[ ! -x $INPUT_BINARY ]] ; then
    echo "INPUT_BINARY must be set to an executable file"
    exit 1
fi

if [[ -z "$PRE_COMMAND" ]] ; then
    PRE_COMMAND=':'
fi

if [[ -z "$POST_COMMAND" ]] ; then
    POST_COMMAND='cat'
fi

if [[ -n "$TIMEOUT_OR_CMD" && $TIMEOUT_OR_CMD =~ ^[0-9]+[smhd]?$ ]] ; then
    TIMEOUT_OR_CMD="timeout -s KILL $TIMEOUT_OR_CMD"
fi

if [[ -z "$MAX_FUNCS_FLAG" ]] ; then
    MAX_FUNCS_FLAG="-max-funcs"
fi

OPTIMIZED_BINARY=$(mktemp -t -u --suffix=.bolt $(basename ${INPUT_BINARY}).XXX)
OUTPUT_FILE="${OPTIMIZED_BINARY}.out"
BOLT_LOG=$(mktemp -t -u --suffix=.log boltXXX)

if [[ -z $OFFLINE ]]; then
    echo "Verify input binary passes"
    echo "  INPUT_BINARY: $PRE_COMMAND && $TIMEOUT_OR_CMD $INPUT_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE"
    ($PRE_COMMAND && $TIMEOUT_OR_CMD $INPUT_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE)
    STATUS=$?
    if [[ "$IGNORE_ERROR" == "1" ]]; then
        FAIL=0
    else
        FAIL=$STATUS
    fi
    if [[ -e "$GOLD_FILE" ]] ; then
        cmp -s "$OUTPUT_FILE" "$GOLD_FILE"
        FAIL=$?
    fi
    if [[ $FAIL -ne "0" ]] ; then
        echo "  Warning: input binary failed"
    else
        echo "  Input binary passes."
    fi
fi

echo "Verify optimized binary fails"
($BOLT $BOLT_OPTIONS $INPUT_BINARY -o $OPTIMIZED_BINARY >& $BOLT_LOG)
FAIL=$?
if [[ $FAIL -eq "0" ]]; then
    if [[ -z $OFFLINE ]]; then
        echo "  OPTIMIZED_BINARY: $PRE_COMMAND && $TIMEOUT_OR_CMD $OPTIMIZED_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE"
        ($PRE_COMMAND && $TIMEOUT_OR_CMD $OPTIMIZED_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE)
        STATUS=$?
        if [[ "$IGNORE_ERROR" == "1" ]]; then
            FAIL=0
        else
            FAIL=$STATUS
        fi
        if [[ -e "$GOLD_FILE" ]] ; then
            cmp -s "$OUTPUT_FILE" "$GOLD_FILE"
            FAIL=$?
        fi
    else
        echo "Did it pass? Type the return code [0 = pass, 1 = fail]"
        read -n1 PASS
    fi
    if [[ $FAIL -eq "0" ]] ; then
        echo "  Warning: optimized binary passes."
    else
        echo "  Optimized binary fails as expected."
    fi
else
    echo "  Bolt crashes while generating optimized binary."
fi

# Collect function names
FF=$(mktemp -t -u --suffix=.txt func-names.XXX)
nm --defined-only -p $INPUT_BINARY | grep " [TtWw] " | cut -d ' ' -f 3 | egrep -v "\._" | egrep -v '^$' | sort -u > $FF

# Use function names or numbers
if [[ -z "$MAX_FUNCS" ]] ; then
    # Do binary search on function names
    if [[ -n "$FUNC_NAMES" ]]; then
        FF=$FUNC_NAMES
    fi
    NUM_FUNCS=$(wc -l $FF | cut -d ' ' -f 1)
    HALF=$(expr \( $NUM_FUNCS + 1 \) / 2)
    PREFIX=$(mktemp -t -u --suffix=.txt func-names.XXX)
    FF0=$PREFIX\aa
    FF1=$PREFIX\ab
    split -a 2 -l $HALF $FF $PREFIX
    FF=$FF0
    NUM_FUNCS=$(wc -l $FF | cut -d ' ' -f 1)
    CONTINUE=$(expr $NUM_FUNCS \> 0)
else
    P=0
    if [[ "$MAX_FUNCS" -eq "0" ]]; then
        Q=$(wc -l $FF | cut -d ' ' -f 1)
    else
        Q=$MAX_FUNCS
    fi
    I=$Q
    CONTINUE=$(expr \( $Q - $P \) \> 1)
fi

ITER=0
while [[ "$CONTINUE" -ne "0" ]] ; do
    rm -f $OPTIMIZED_BINARY
    if [[ -z "$MAX_FUNCS" ]] ; then
        echo "Iteration $ITER, trying $FF / $HALF functions"
        SEARCH_OPT="-funcs-file-no-regex=$FF"
    else
        I=$(expr \( $Q + $P \) / 2)
        echo "Iteration $ITER, P=$P, Q=$Q, I=$I"
        SEARCH_OPT="$MAX_FUNCS_FLAG=$I"
    fi
    echo "  BOLT: $BOLT $BOLT_OPTIONS $INPUT_BINARY $SEARCH_OPT -o $OPTIMIZED_BINARY >& $BOLT_LOG"
    ($BOLT $BOLT_OPTIONS $INPUT_BINARY $SEARCH_OPT -o $OPTIMIZED_BINARY >& $BOLT_LOG)
    FAIL=$?
    echo "  BOLT failure=$FAIL"
    rm -f $OUTPUT_FILE
    if [[ $FAIL -eq "0" ]] ; then
        if [[ -z $OFFLINE ]]; then
            echo "  OPTIMIZED_BINARY: $PRE_COMMAND && $TIMEOUT_OR_CMD $OPTIMIZED_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE"
            ($PRE_COMMAND && $TIMEOUT_OR_CMD $OPTIMIZED_BINARY $COMMAND_LINE |& $POST_COMMAND >& $OUTPUT_FILE)
            STATUS=$?
            if [[ "$IGNORE_ERROR" == "1" ]]; then
                FAIL=0
            else
                FAIL=$STATUS
            fi
            if [[ -e "$GOLD_FILE" ]] ; then
                cmp -s "$OUTPUT_FILE" "$GOLD_FILE"
                FAIL=$?
            fi
            echo "  OPTIMIZED_BINARY failure=$FAIL"
        else
            echo "Did it pass? Type the return code [0 = pass, 1 = fail]"
            read -n1 PASS
        fi
    else
        FAIL=1
    fi

    if [[ -z "$MAX_FUNCS" ]] ; then
        if [[ $FAIL -eq "0" ]] ; then
            if [[ "$FF" == "$FF1" ]]; then
                NUM_FUNCS=0
                break;
            fi
            FF=$FF1
            NUM_FUNCS=$(wc -l $FF | cut -d ' ' -f 1)
        else
            HALF=$(expr \( $NUM_FUNCS + 1 \) / 2)
            PREFIX=$(mktemp -t -u --suffix=.txt func-names.XXX)
            split -a 2 -l $HALF $FF $PREFIX
            FF0=$PREFIX\aa
            FF1=$PREFIX\ab
            FF=$FF0
            NUM_FUNCS=$(wc -l $FF | cut -d ' ' -f 1)
            if [[ $NUM_FUNCS -eq "1" && ! -e $FF1 ]]; then
                break;
            fi
        fi
        CONTINUE=$(expr $NUM_FUNCS \> 0)
    else
        if [[ $FAIL -eq "0" ]] ; then
            P=$I
        else
            Q=$I
        fi
        FF=$I
        HALF=$I
        CONTINUE=$(expr \( $Q - $P \) \> 1)
    fi
    ITER=$(expr $ITER + 1)
done

if [[ -z "$MAX_FUNCS" ]] ; then
    if [[ "$NUM_FUNCS" -ne "0" ]] ; then
        FAILED="The function(s) that failed are in $FF"
    fi
else
    if [[ $P -ne $Q ]] ; then
        FF=$(grep "processing ending" $BOLT_LOG | sed -e "s/BOLT-INFO: processing ending on \(.*\)/\1/g" | tail -1)
        FAILED="The item that failed is $FF @ $Q"
    fi
fi

if [[ -n "$FAILED" ]] ; then
    echo "$FAILED"
    echo "To reproduce, run: $BOLT $BOLT_OPTIONS $INPUT_BINARY $SEARCH_OPT -o $OPTIMIZED_BINARY"
else
    echo "Unable to reproduce bug."
fi

rm $OPTIMIZED_BINARY $OUTPUT_FILE $BOLT_LOG