let size=readNumber("Enter size of array:");
let a=makeArray(size);
println("Made array of length:"+len(a));
i=0;
loop i<len(a) {
    a[i]=readNumber("Enter element "+i+":");
    // a[i]=readNumber("Enter element "+i);
    ++i;
}
println(a);
i=0;
loop i<len(a) {
    println(a[i]);
    ++i;
}