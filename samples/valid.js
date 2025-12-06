// incluir "stdlib"

funcion principal() int {
    definir int kk = (5*6-4)+(2*3)
    definir int a = (5+8)*doubled(2), b= 0xDEADBEEF
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