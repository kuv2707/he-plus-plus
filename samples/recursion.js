recprint(10);
function recprint(a){
    if (a==0) {
        return 0;
    }
    recprint(a-1);
    print(a+",");
    return 0;
}
