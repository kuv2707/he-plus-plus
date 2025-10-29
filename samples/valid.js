// incluir "stdlib"

funcion principal() int {
    definir int a = (5+8)*doubled(2)
    // definir int b = [int]{1,2,3}[1]
    // definir int c = a + b + (2+3)*5
    // definir int d = c
    printf("%d", a)
    devolver a
}

funcion doubled(a int) int {
    devolver 2*a
}