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