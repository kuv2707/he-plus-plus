funcion principal() vacio {
    definir int c = 15
    definir int d = 15
    definir int e = 15
    definir int f = d + e
    f = f + 1

    definir [int] arr = [int][c]
    arr[0] = 0
    arr[1] = 1
    para definir int i=2;i<c;i=i+1 {
        arr[i] = arr[i-1] + arr[i-2]
        log(arr[i])
        // e = e + 1
        // definir int k = i + 1
        // doubled(e, k)
    }
}

funcion doubled(a int, b int) int {
    devolver 2*a
}

funcion log(c int) vacio {

}