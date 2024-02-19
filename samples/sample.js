function abs(arg){
        if arg>=0 {
            return arg;
        }else{
            return -arg;
        }
    }
arr=[1,2,3,-4,-5,-6];
i=0;
loop i<len(arr){
    print(abs(arr[i]));
    i=i+1;
}