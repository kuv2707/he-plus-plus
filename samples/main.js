funcion principal() &[[int]] {
    // hello()()
    definir int kkkk = 5+verdad
    definir int a = 11
    definir float fpt = 5.76
    definir int expr = 5*6/(3-4-a)
    definir &int ptr = &a
    definir [inta] arr = [int]{1+2, 2, 3-2}
    // definir bool bvar = 1 == 2
    println(*ptr)
    perform(3, 4, subtract)
    // devolver ptr
    devolver &[[int]]{arr,[int]{7,6}}
}

funcion add(a int, b int) int {
    devolver a+b
}

funcion subtract(a int, b int) int {
    devolver a+b
}

funcion perform(aa int, bb int, op funcion(int,int)int) int {
    op(aa,bb,4)
}

funcion println(a int) vacio {
    printf("Hello %d", a)
}