function give2(){
    return 2+give3();
}
function give3(){
    return 3+give4();
}
function give4(){
    return 4;
}
println(give2());