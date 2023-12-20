//find k-root of number n by bisection method
findroot(27,3);

func(n,k)findroot{
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
    print(high);
    
    func(arg)abs{
        if arg>=0 {
            return arg;
        }else{
            return 0-arg;
        }
    }
    func(😶‍🌫️)print{
        #(😶‍🌫️);
    }
    func(a,n)pow{
        let ret=1;
        loop n>0 {
            ret=ret*a;
            n=n-1;
        }
        return ret;
    }
}