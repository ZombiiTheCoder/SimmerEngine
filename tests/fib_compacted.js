function fib(num){var a=1;var b=0;var temp=null;while (num >= 0){temp = a;a = a + b;b = temp;num--};};fib(1000)