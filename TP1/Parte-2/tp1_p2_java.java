import java.time.Duration;
import java.time.Instant;
import java.util.Arrays;
import java.util.Random;

public class tp1_p2_java {
  private static final int ORDER = 9000;
  private static final int THREADSNUM = 100;
  private static final int SCALAR = 7;
  private static final int[][] mat = new int[ORDER][ORDER];
  private static final int[][] matConc = new int[ORDER][ORDER];
  private static final int[][] matSec = new int[ORDER][ORDER];
  private static final Random rand = new Random();

  public static void main(final String[] args) {

    initializeMatrix();

    System.out.println("scalar: " + SCALAR + " order: " + ORDER);

    if (ORDER <= 10) {
      printMatrix(mat);
    }

    final Instant startRs = Instant.now();
    secuentialScalarMultiplication();
    final Instant finishRs = Instant.now();

    System.out.println("Time of Sequential Matrix: " + Duration.between(startRs, finishRs).toMillis() + " ms");

    if (ORDER <= 10) {
      System.out.println("Matriz secuencial");
      printMatrix(matSec);
    }

    final Instant startRc = Instant.now();
    concurrentScalarMultiplication();
    final Instant finishRc = Instant.now();

    System.out.println("Time of Concurrent Matrix: " + Duration.between(startRc, finishRc).toMillis() + " ms");
    System.out.println("Using " + THREADSNUM + " threads");

    if (ORDER <= 10) {
      System.out.println("Matriz concurrente");
      printMatrix(matConc);
    }

    System.out.println("Are matrices matConc and matSec equal? " + Arrays.deepEquals(matConc, matSec));
  }

  public static void secuentialScalarMultiplication() {
    for (int i = 0; i < ORDER; i++) {
      for (int j = 0; j < ORDER; j++) {
        matSec[i][j] = mat[i][j] * SCALAR;
      }
    }
  }

  public static void concurrentScalarMultiplication() {
    final Thread[] threads = new Thread[THREADSNUM];

    final int rowsPerThread = ORDER / THREADSNUM;

    for (int i = 0; i < THREADSNUM; i++) {
      final int startRow = i * rowsPerThread;
      // Esta línea garantiza que el último hilo procese todas las filas restantes si
      // el orden no es perfectamente divisible por THREADSNUM.
      // Garantiza que se procesen todas las filas, distribuyendo el trabajo de manera
      // más uniforme.
      final int endRow = (i == THREADSNUM - 1) ? ORDER : (i + 1) * rowsPerThread;

      threads[i] = new Thread(() -> {
        for (int k = startRow; k < endRow; k++) {
          for (int w = 0; w < ORDER; w++) {
            matConc[k][w] = mat[k][w] * SCALAR;
          }
        }
      });

      threads[i].start();

    }

    for (final Thread thread : threads) {
      try {
        thread.join();
      } catch (final InterruptedException e) {
        e.printStackTrace();
        Thread.currentThread().interrupt();
      }
    }
  }

  public static void printMatrix(final int[][] mat) {
    for (final int[] row : mat) {
      for (final int col : row) {
        System.out.print(col + " ");
      }
      System.out.println();
    }

  }

  public static void initializeMatrix() {
    final int min = 0;
    final int max = 9;

    for (int i = 0; i < ORDER; i++) {
      for (int j = 0; j < ORDER; j++) {
        mat[i][j] = rand.nextInt(max - min + 1) + min;
      }
    }

  }
}
