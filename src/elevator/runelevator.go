package elevator

import(
		."fmt"
		"time"
		//.".././network"
		.".././channels"
		)

const N_BUTTONS int = 3
const N_FLOORS int = 4

type elev_button_type_t int
const(
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
	)
type elev_motor_direction_t int
const(
	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP = 1
	)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

var dirn int

func elev_init() bool{ 
	if io_init() == 0{ //initialization of i/o
		return false
	}
	for i:=0;i<N_FLOORS;i++{
		if i!=0{
			elev_set_button_lamp(BUTTON_CALL_DOWN, i, false) //set all call down buttons to 0
		}
		if i!=N_FLOORS-1{
			elev_set_button_lamp(BUTTON_CALL_UP, i, false)	//set all call up buttons to 0
		}		
		elev_set_button_lamp(BUTTON_COMMAND, i, false) //set all floor buttons to 0
	}
	elev_set_stop_lamp(false)
	elev_set_door_open_lamp(false)
	elev_set_floor_indicator(0)
	if elev_get_floor_sensor_signal() != 0{
		elev_set_motor_direction(DIRN_DOWN)
		for{
			if elev_get_floor_sensor_signal() != -1{
				elev_set_motor_direction(DIRN_UP)
				time.Sleep(10*time.Millisecond)					
				elev_set_motor_direction(DIRN_STOP)
				dirn = DIRN_STOP
				break
			}
		}			
	}

	return true
}

func elev_set_button_lamp(button elev_button_type_t, floor int, value bool){
	if value == true{
		io_set_bit(lamp_channel_matrix[floor][button])
	}
	if value == false{
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}

func elev_set_stop_lamp(value bool){
	if value == true{
		io_set_bit(LIGHT_STOP)
	}else{
		io_clear_bit(LIGHT_STOP)
	}
}

func elev_set_door_open_lamp(value bool){
	if value == true{
		io_set_bit(LIGHT_DOOR_OPEN)
	}else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}	
}
func elev_set_floor_indicator(floor int){
	switch floor{
	case 0:
		io_clear_bit(LIGHT_FLOOR_IND1)
		io_clear_bit(LIGHT_FLOOR_IND2)
	case 1:
		io_clear_bit(LIGHT_FLOOR_IND1)
		io_set_bit(LIGHT_FLOOR_IND2)
	case 2:
		io_set_bit(LIGHT_FLOOR_IND1)
		io_clear_bit(LIGHT_FLOOR_IND2)
	case 3:
		io_set_bit(LIGHT_FLOOR_IND1)
		io_set_bit(LIGHT_FLOOR_IND2)			
	}
}

func elev_set_motor_direction(dirn elev_motor_direction_t){
	if dirn == 0{
		io_write_analog(MOTOR, 0)
	}else if dirn > 0{
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR,2800)
	}else if dirn < 0{
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR,2800)	
	}
}

func elev_get_floor_sensor_signal() int{
	if io_read_bit(SENSOR_FLOOR1) == 1{
		return 0
	}else if io_read_bit(SENSOR_FLOOR2) == 1{
		return 1
	}else if io_read_bit(SENSOR_FLOOR3) == 1{
		return 2
	}else if io_read_bit(SENSOR_FLOOR4) == 1{
		return 3
	}
	return -1	 
}

func elev_get_button_signal(button elev_button_type_t, floor int) bool{
	if io_read_bit(button_channel_matrix[floor][button]) != 0{
		return true	
	}else{
		return false
	}
} 

func elev_set_button_signal(button elev_button_type_t, floor int, value bool){
	if value == true{
		io_set_bit(lamp_channel_matrix[floor][button])
	}else{
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}

func elev_get_stop_signal() int{
	return io_read_bit(STOP)
}

func newInternalOrders(internalOrders [N_FLOORS]int) [N_FLOORS]int{
	for i:=0;i<N_FLOORS;i++{
		if elev_get_button_signal(BUTTON_COMMAND,i) == true{
			internalOrders[i] = 1	
		} 
	}
	return internalOrders
}
func clearInternalOrders(dirn int,floor int, internalOrders [N_FLOORS]int) [N_FLOORS]int{
	if floor != -1 && dirn == 0{
		internalOrders[floor] = 0
	}
	return internalOrders
}

func setInternalLights(internalOrders [N_FLOORS]int){
	for i:=0;i < N_FLOORS;i++{
		if internalOrders[i] == 1{
			elev_set_button_signal(BUTTON_COMMAND,i,true)	
		}else{
			elev_set_button_signal(BUTTON_COMMAND,i,false)
		}
	}
}

func newExternalOrders(externalOrders [N_FLOORS][2]int) [N_FLOORS][2]int{
	for i:=0;i<N_FLOORS-1;i++{
		if elev_get_button_signal(BUTTON_CALL_UP,i) == true{
			externalOrders[i][0] = 1
		}
	}
	for i:=1;i<N_FLOORS;i++{
		if elev_get_button_signal(BUTTON_CALL_DOWN,i) == true{
			externalOrders[i][1] = 1
		}
	}	
	return externalOrders
}

func clearExternalOrders(dirn int, floor int, externalOrders [N_FLOORS][2]int) [N_FLOORS][2]int{
	if floor != -1 && dirn == 0{
		externalOrders[floor][0] = 0
		externalOrders[floor][1] = 0
	}
	return externalOrders
}

func setExternalLights(externalOrders [N_FLOORS][2]int){
	for i:=0;i<N_FLOORS-1;i++{
		if externalOrders[i][0] == 1{
			elev_set_button_lamp(BUTTON_CALL_UP,i,true)
		}else{
			elev_set_button_lamp(BUTTON_CALL_UP,i,false)
		}
	}
	for i:=1;i<N_FLOORS;i++{
		if externalOrders[i][1] == 1{
			elev_set_button_lamp(BUTTON_CALL_DOWN,i,true)
		}else{
			elev_set_button_lamp(BUTTON_CALL_DOWN,i,false)
		}
	}
}

func InformationToNetworkUnit(internalOrders [N_FLOORS]int,externalOrders [N_FLOORS][2]int, ExternalOrdersToNetwork chan [N_FLOORS][2]int, InternalOrdersToNetwork chan [N_FLOORS]int){
	InternalOrdersToNetwork <- internalOrders
	ExternalOrdersToNetwork <- externalOrders
	//Direction <- dirn
}
func States(FloorChan chan int){
	lastFloor := elev_get_floor_sensor_signal()
	for{
			if elev_get_floor_sensor_signal() != -1{
				lastFloor = elev_get_floor_sensor_signal()
			}					

			FloorChan <- lastFloor

	}	
}

func stopAtFloor(order []int,lastFloor int)bool{
	stop := false
	for i:= 0; i < len(order);i++{
		if order[i] == 1{
			if i == lastFloor{
				stop = true
			}	
		}
	} 
	return stop
}


func ExecuteOrders(ExecuteListChan chan []int, DirectionChan chan int){
	go States(FloorChan)
	lastFloor := elev_get_floor_sensor_signal()
	for{
		direction :=<- DirectionChan
		order :=<- ExecuteListChan				
		switch{
		case elev_get_floor_sensor_signal() != -1: //elevator is at a floor
			lastFloor = elev_get_floor_sensor_signal()
			if stopAtFloor(order,lastFloor) == true{
				elev_set_motor_direction(DIRN_STOP)
				DirnChan <- 0
				/*set door open
				for i:= 0;i<50;i++{ 
					DirnChan <- 0
					<- DirectionChan
					<- ExecuteListChan
					time.Sleep(25*time.Millisecond)
				}
				clear doo open*/
			}
			if direction > 0{
				elev_set_motor_direction(DIRN_UP)
			}
			if direction < 0{
				elev_set_motor_direction(DIRN_DOWN)
			}
			if direction == 0{
				elev_set_motor_direction(DIRN_STOP)
			}

		case elev_get_floor_sensor_signal() == -1:			
					
		}
	}
}

/*
func ExecuteOrdersFromNetworkUnit(ExecuteListChan chan []int, FloorChan chan int,NextDirChan chan int){
	lastFloor := elev_get_floor_sensor_signal()
	nextDir := DIRN_STOP
	go States(ExecuteListChan,FloorChan,NextDirChan)
	for{
		executeFromNetwork :=<- ExecuteListChan
		Println("yolo")
		Println(executeFromNetwork)
		
		select{
		case executeFromNetwork :=<- ExecuteListChan:

			if elev_get_floor_sensor_signal() != -1{ //sets current floor
				lastFloor = elev_get_floor_sensor_signal()

			}		
			if lastFloor == executeFromNetwork[0]{ //if right floor; stop
				elev_set_motor_direction(DIRN_STOP)
			}

			if lastFloor < executeFromNetwork[0] && executeFromNetwork[0] != -1{ //if order is 

				elev_set_motor_direction(DIRN_UP)
				for{
					
					if elev_get_floor_sensor_signal() != -1{
						lastFloor = elev_get_floor_sensor_signal()
					}					
					if elev_get_floor_sensor_signal() == executeFromNetwork[0]{
						elev_set_motor_direction(DIRN_DOWN)
						time.Sleep(10*time.Millisecond)					
						elev_set_motor_direction(DIRN_STOP)
						//set door open lamp
						time.Sleep(100*time.Millisecond)
						//clear door open lamp

						break
					}
				}
			
			}

			if lastFloor > executeFromNetwork[0] && executeFromNetwork[0] != -1{
				elev_set_motor_direction(DIRN_DOWN)
				for{

					if elev_get_floor_sensor_signal() != -1{
						lastFloor = elev_get_floor_sensor_signal()
					}					
					if elev_get_floor_sensor_signal() == executeFromNetwork[0]{
						elev_set_motor_direction(DIRN_UP)
						time.Sleep(10*time.Millisecond)					
						elev_set_motor_direction(DIRN_STOP)
						//set door open lamp
						time.Sleep(100*time.Millisecond)
						//clear door open lamp						
		
						break
					}
				}			
			}
		default:
			elev_set_motor_direction(DIRN_STOP)					
		}
	}
}
*/

func lightsAndOrders(internalOrders [N_FLOORS]int, externalOrders [N_FLOORS][2]int, DirnChan chan int){	
	for{
		select{
		case dirn :=<-DirnChan:
			internalOrders = clearInternalOrders(dirn,elev_get_floor_sensor_signal(),internalOrders)
			externalOrders = clearExternalOrders(dirn,elev_get_floor_sensor_signal(),externalOrders)
			default:
			internalOrders = newInternalOrders(internalOrders)
			externalOrders = newExternalOrders(externalOrders)
			setInternalLights(internalOrders)
			setExternalLights(externalOrders)
			elev_set_floor_indicator(elev_get_floor_sensor_signal())
			InformationToNetworkUnit(internalOrders,externalOrders,ExternalOrdersToNetwork,InternalOrdersToNetwork/*,dirn,Direction*/)			
	
		}
	}
	
}

func runElevator(){
	elev_init()
	elev_set_motor_direction(DIRN_STOP)
	//go ExecuteOrdersFromNetworkUnit(ExecuteListChan,DirnChan,dirn,FloorChan,NextDirChan)		
	go ExecuteOrders(ExecuteListChan,DirectionChan)
}

func Elevator(){
	internalOrders := [N_FLOORS]int{0,0,0,0}
	externalOrders := [N_FLOORS][2]int{{0,0},{0,0},{0,0},{0,0}} //{Up,Down}
	
	
	go runElevator()
	go lightsAndOrders(internalOrders,externalOrders,DirnChan)
	

	Println("Elevator")
}

