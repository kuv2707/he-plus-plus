funcion principal() vacio {
    definir int c = 15*(4+1)
  

    definir [int] arr = [int][c]
    arr[0] = 0
    arr[1] = 1
    para definir int i=2;i<c;i=i+1 {
        arr[i] = add(arr[i-1], arr[i-2])
        
    }
}

funcion add(a int, b int) int {
    devolver a+b
}

funcion log(c int) int {
    devolver c
}