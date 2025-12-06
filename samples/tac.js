funcion principal() vacio {
    definir int a = 502
    definir int b = 6
    definir int c = a + doubled(a+b, a)
    si a < b entonces {
        a = a + 5
    } o si a == b entonces {
        a = a
    } o {
        a = a - 5
    }
}

funcion doubled(a int, b int) int {
    devolver 2*a
}