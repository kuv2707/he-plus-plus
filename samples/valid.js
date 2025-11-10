// incluir "stdlib"

funcion principal() int {
    definir int a = (5+8)*doubled(2)
    definir int b = 0xDEADBEEF
    definir int c = a + b + (2+3)*5
    si a<c entonces {
        printf("greater")
    }

    definir int d = a
    printf("%d", a)
    devolver a
}

funcion doubled(a int) int {
    devolver 2*a
}