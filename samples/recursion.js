func(a)recprint{
    if a==0 {
        print(a);
        return 0;
    }
    recprint(a-1);
    print(a);
    return 0;
}
func(a)print{
    #(a);
}

recprint(10);