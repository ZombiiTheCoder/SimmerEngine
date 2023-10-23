function fib(num){
    var a = 1, b = 0, temp, i = 0;
    while (i < num){
      temp = a;
      a = a + b;
      b = temp;
      i++;
    }
  
    return b;
}

fib(12)
"223123"