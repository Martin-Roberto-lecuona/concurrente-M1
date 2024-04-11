import java.io.IOException;
import java.util.ArrayList;
import java.util.Map;
import java.util.Set;

class ProcessTree
{
  static Map<Character, Set<Character>> processMap = Map.of (
      'A', Set.of('B', 'C', 'D'),
      'B', Set.of('E', 'F'),
      'F', Set.of('H', 'I'),
      'I', Set.of('J'),
      'D', Set.of('G')
  );

  static ArrayList<Process> processList = new ArrayList<>();

  static final int SLEEP_TIME = 100000;
  // Milliseconds

  public static void showInfo ( char processName )
  {
    ProcessHandle me = ProcessHandle.current();
    long pid  = me.pid();
    long ppid = me.parent().get().pid();
    System.out.println( "Process " + processName + ", PID: " + pid
        + ", PPID: " + ppid );
  }

  public static void createChilds ( char processName ) throws IOException,
      InterruptedException
  {
    if ( processMap.containsKey(processName) )
    {
      for ( char son : processMap.get(processName) )
      {
        ProcessBuilder builder = new ProcessBuilder("java", "ProcessTree.java",
            String.valueOf(son));
        builder.inheritIO();
        Process process = builder.start();
        processList.add(process);
      }
    }
    else
    {
      Thread.sleep(SLEEP_TIME);
    }
  }

  public static void waitChilds ( ) throws InterruptedException
  {
    for ( Process p : processList )
    {
      p.waitFor();
    }
  }

  public static void main ( String[] args ) throws IOException,
      InterruptedException
  {
    char root = 'A';
    char process = ( args.length > 0 ) ? args[0].charAt(0) : root;
    showInfo(process);
    createChilds(process);
    waitChilds();
  }
}