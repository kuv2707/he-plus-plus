//find k-root of number n by bisection method
// findroot(readNumber("Enter a number:"),2);
findroot(1600,2);
function findroot (n,k){
    // println(n);
    // println(k);
    let low=0;
    let high=n;
    let mid=0;    
    let TOLERANCE=0.000000001;
    loop abs(low-high)>TOLERANCE {
        mid=(low+high)/2;
        if pow(mid,k)>n {
            high=mid;
        }else{
            low=mid;
        }
    }
    print("The square root is:");
    println(high);
    
    function abs(arg){
        if arg>=0 {
            return arg;
        }else{
            return -arg;
        }
    }
    function pow(a,n){
        let ret=1;
        loop n>0 {
            ret=ret*a;
            n=n-1;
        }
        return ret;
    }
}