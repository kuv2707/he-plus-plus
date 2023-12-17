//find square root by bisection method
let n=9;
let low=0;
let high=n;
let mid=0;    
let TOLERANCE=0.00001;
loop(abs(low-high)>TOLERANCE){
    mid=(low+high)/2;
    if(mid*mid>n){
        high=mid;
    }else{
        low=mid;
    }
}
#(low/2+high/2);

func(arg)abs{
    if(arg>=0){
        return arg;
    }else{
        return 0-arg;
    }
}