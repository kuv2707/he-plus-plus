println("Hello "+five(5));

function five(a){
    if(a == 0) {
        return 1;
    }
    return 5*five(a-1);
}