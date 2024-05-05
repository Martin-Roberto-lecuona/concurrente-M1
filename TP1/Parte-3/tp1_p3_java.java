import java.util.NoSuchElementException;
import java.util.Scanner;

public class Parte3 
{
    static final String SIDE_A = "LADO A";
    static final String SIDE_B = "LADO B";
    static final String BEGIN = SIDE_A;
    static final String FOX = "Zorro";
    static final String CHICKEN = "Pollo";
    static final String MAIZ = "Maiz";
    static final String FARMER = "Granjero";
    static final int OPTION_FARMER = 0;
    static final int OPTION_FOX = 1;
    static final int OPTION_CHICKEN = 2;
    static final int OPTION_MAIZ = 3;
    public static void main(String[] args) 
    {
        InitialSide sideA = new InitialSide(SIDE_A);
        EndSide sideB = new EndSide(SIDE_B);

        sideA.setOtherSide(sideB);
        sideB.setOtherSide(sideA);

        sideA.start();
        sideB.start();
    }
    static public class SideThread extends Thread 
    {
        protected String name;
        public StateLado state;
        public volatile boolean running = true;
        private String lastObject;

        public SideThread(String name) 
        {
            this.name = name;
        }
        public synchronized void stopThread() 
        {
            this.running = false;
        }
    
        public String getNameSide() 
        {
            return this.name;
        }
    
        public synchronized boolean isRunning() 
        {
            return running;
        }

        public String getLastObject()
        {
            return this.lastObject;
        }
        public void setLastObject(String object)
        {
            this.lastObject = object;
        }
    }
    static class InitialSide extends SideThread 
    {
        private EndSide otherSide;

        public InitialSide(String name) 
        {
            super(name);
            this.state = new StateLado(name,true);
        }
        public void setOtherSide(EndSide otherSide) 
        {
            this.otherSide = otherSide;
        }

        public synchronized void run() 
        {
            Scanner scanner = new Scanner(System.in);    
            try{
                Mesagges.welcome();                 
                while(running)
                {
                    this.state.menuBack(scanner,otherSide);
                    if(!this.state.validate())
                    {
                        running = false;
                        otherSide.stopThread();
                        synchronized (otherSide)
                        {
                            otherSide.notify();
                        }
                        continue;
                    }
                    synchronized (otherSide) 
                    {
                        otherSide.notify();
                    }
                    wait();
                    if(!running)
                        break;
                    if(this.getLastObject() != null)
                        this.state.refresh_states(this.getLastObject());
                    else
                        Mesagges.farmerAlone(this.state.site);

                    this.state.farmer_in_site = !this.state.farmer_in_site;
                    this.setLastObject(null);
                }
            
            }catch (InterruptedException e) 
            {
                e.printStackTrace();
            }
            catch (NoSuchElementException e) 
            {
                Mesagges.errorScanner();
            }finally
            {
                scanner.close();
                running = false;
                otherSide.stopThread();
            }
        }
           
    }

    static class EndSide extends SideThread {
        private InitialSide otherSide;

        public EndSide(String name) {
            super(name);
            this.state = new StateLado(name);
        }
        public void setOtherSide(InitialSide otherSide) {
            this.otherSide = otherSide;
        }

        public synchronized void run() {
            Scanner scanner = new Scanner(System.in);  
            while (running) {
                try {
                    synchronized (this) {
                        wait();
                    }
                    if(!running)
                        break;
                    String last = this.getLastObject();
                    if(last != null && !last.isEmpty())
                    {
                        this.state.refresh_states(last);
                    }
                    else
                    {
                        Mesagges.farmerAlone(this.name);
                    }
                    this.state.farmer_in_site = !this.state.farmer_in_site;
                    if (this.state.checkWin())
                    {
                        Mesagges.winRiddle();
                        running = false;
                        otherSide.stopThread();
                        synchronized (otherSide)
                        {
                            otherSide.notify();
                        }
                        break;
                    }
                    state.menuBack(scanner,otherSide);
                    if(!this.state.validate())
                    {
                        running = false;
                        otherSide.stopThread();
                    }
                    this.setLastObject(null);
                    synchronized (otherSide) {
                        Mesagges.crossChange(name,otherSide.name);
                        otherSide.notify();
                    }
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }
    }
    
    static class StateLado
    {
        private boolean fox_in_site;
        private boolean maiz_in_site;
        private boolean chicken_in_site;
        private boolean farmer_in_site;
        private String site;

        public StateLado(String site)
        {
            this.fox_in_site = false;
            this.maiz_in_site = false;
            this.chicken_in_site = false;
            this.farmer_in_site = false;
            this.site = site;
        }
        public StateLado(String site,boolean begin)
        {
            this.fox_in_site = begin;
            this.maiz_in_site = begin;
            this.chicken_in_site = begin;
            this.farmer_in_site = begin;
            this.site = site;
        }
        public boolean getFoxInSite(){
            return this.fox_in_site;
        }
        public void setFoxInSite(boolean fox_in_site){
            this.fox_in_site = fox_in_site;
        }
        public boolean getMaizInSite(){
            return this.maiz_in_site;
        }
        public void setMaizInSite(boolean maiz_in_site){
            this.maiz_in_site = maiz_in_site;
        }
        public boolean getChickenInSite(){
            return this.chicken_in_site;
        }
        public void setChickenInSite(boolean chicken_in_site){
            this.chicken_in_site = chicken_in_site;
        }
        public boolean getFarmerInSite(){
            return this.farmer_in_site;
        }
        public void setFarmerInSite(boolean farmer_in_site){
            this.farmer_in_site = farmer_in_site;
        }

        public synchronized boolean checkWin()
        {
            if(this.chicken_in_site && this.fox_in_site && this.maiz_in_site && this.farmer_in_site)
            {
                return true;
            }
            return false;
        }
        public synchronized void cross(String objeto, String otherSideString) 
        {
            switch (objeto) {
                case FARMER:
                    break;
                case FOX:
                    this.fox_in_site = !this.fox_in_site;
                    break;
                case CHICKEN:
                    this.chicken_in_site = !this.chicken_in_site;
                    break;
                case MAIZ:
                    this.maiz_in_site = !this.maiz_in_site;
                    break;
            }
            this.farmer_in_site = !this.farmer_in_site;
        }
        public synchronized boolean validate() 
        {
            if (this.fox_in_site && this.chicken_in_site && !this.farmer_in_site) 
            {
                Mesagges.messageValidate(FOX, CHICKEN,this.site);
                return false;
            } else if (this.chicken_in_site && this.maiz_in_site && !this.farmer_in_site) 
            {
                Mesagges.messageValidate(CHICKEN, MAIZ, this.site);
                return false;
            }
            return true;
        }
        public void menuBack(Scanner scanner,SideThread otherSide)
        {
            while (true) 
            {
                System.out.println("======================== MENÚ LADO "+ site +" ========================");
                System.out.println("Por favor, ingrese cómo desea que el granjero viaje:");
                System.out.println(this.farmer_in_site ? "0. Granjero solo" : "");
                System.out.println(this.fox_in_site ? "1. Zorro" : "");
                System.out.println(this.chicken_in_site ? "2. Chicken" : "");
                System.out.println(this.maiz_in_site ? "3. Maiz" : "");
                System.out.print("¿Cruzar con?: ");
                int option = scanner.nextInt();
                if(processOption(option,otherSide) == true)
                    break;
            }
        }

        public synchronized boolean processOption(int option, SideThread otherSide) 
        {
            switch (option) 
            {
                case OPTION_FARMER:
                    cross(FARMER, otherSide.getNameSide());
                    break;
                case OPTION_FOX:
                    if (this.fox_in_site) 
                    {
                        cross(FOX, otherSide.getNameSide());
                        otherSide.setLastObject(FOX);
                    } else {
                        Mesagges.messageOtherSide(FOX);
                        return false;
                    }
                    break;
                case OPTION_CHICKEN:
                    if (this.chicken_in_site) 
                    {
                        cross(CHICKEN, otherSide.getNameSide());
                        otherSide.setLastObject(CHICKEN);
                    } else 
                    {
                        Mesagges.messageOtherSide(CHICKEN);
                        return false;
                    }
                    break;
                case OPTION_MAIZ:
                    if (this.maiz_in_site) 
                    {
                        cross(MAIZ, otherSide.getNameSide());
                        otherSide.setLastObject(MAIZ);
                    } else 
                    {
                        Mesagges.messageOtherSide(MAIZ);
                        return false;
                    }
                    break;
                default:
                    Mesagges.messageError();
                    return false;
            }
            return true;
        }
        public void refresh_states(String option)
        {
            if(option != null)
            {
                switch (option) 
                {
                    case FOX:
                        this.fox_in_site = !this.fox_in_site;
                        Mesagges.messageCross(FOX, site);
                        break;
                    case CHICKEN:
                        this.chicken_in_site = !this.chicken_in_site;
                        Mesagges.messageCross(CHICKEN, site);
                        break;
                    case MAIZ:
                        this.maiz_in_site = !this.maiz_in_site;
                        Mesagges.messageCross(MAIZ, site);
                        break;
                    default:
                        break;
                }
            }
            
        }

    }
    public class Mesagges 
    {
        public static void welcome() 
        {
            System.out.println("Bienvenido al problema del Zorro, el Pollo y el Maiz.");
        }
        public static void messageError() 
        {
            System.out.println("Se ha producido un error. Inténtalo de nuevo.");
        }
        public static void invalidOption() 
        {
            System.out.println("Opción inválida.");
        }
        public static void messageOtherSide(String object) 
        {
            System.out.println("El " + object +" ya está del otro lado.");
        }
        public static void messageValidate(String object, String otherObject, String side) 
        {
            System.out.println("¡La "+ object +" se comió el "+ otherObject +"! En el "+side);
            System.out.println("Fin del programa");
        }
        public static void messageCross(String object, String side) 
        {
            if(object==null)
                System.out.println(FARMER + " se encuentran en "+ side);
            else
                System.out.println(FARMER + " y " + object + " se encuentran en "+side);
        }
        public static void crossChange(String side, String otherSide) 
        {
            System.out.println("Cruzando de " + side + " a " + otherSide);
        }
        public static void farmerAlone(String side) 
        {
            System.out.println("El granjero llego a " + side + " solo." );
        }
        
        public static void errorScanner() 
        {
            System.out.println("Se ha producido un error al leer la entrada. Inténtalo de nuevo.");
        }
        public static void winRiddle()
        {
            System.out.println("Todos se encuentran en lado B");
            System.out.println("Fin del programa");
        }
    }
}

