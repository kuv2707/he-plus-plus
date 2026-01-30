funcion principal() vacio {
    definir int c = 15*(4+1)
  

    definir [int] arr = [int][c]
    arr[0] = 0
    arr[1] = 1
    para definir int i=0;i<c;i=i+1 {
        si i > 1 entonces {
            // arr[i] = (arr[i-1] + arr[i-2])
            arr[i] = arr[i-1]
        }
        
    }
}

funcion add(a int, b int) int {
    devolver a+b
}

funcion log(c int) int {
    definir int a = 5
    definir int b = 6
    log(a+b-c)
    devolver c
}