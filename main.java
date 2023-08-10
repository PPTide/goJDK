// Type your code here, or load an example.
class Square {
    static int square(int num) {
        return num * num;
    }
    static int times2(int num) {
        return num * 2;
    }
    public static void main(String[] args) {
        int eight = times2(4);
        System.out.println(square(eight));
    }
}

class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hi 3^2 = " + Square.square(3));
    }
}

class Hello {
    public static int fib(int n) {
        if (n <= 1)
            return n;
        return fib(n - 1) + fib(n - 2);
    }

    public static void main(String args[]) {
        int N = 100;
        for (int i = 0; i < N; i++) {
            System.out.print(fib(i) + " ");
        }
    }
}