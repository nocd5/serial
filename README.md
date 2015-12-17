# serial
CLI tool for serial port connection

## Usage
`serial -p COM3 -b 9600`

## Options

### Required
* -p / --port  
  Port Name
* -b / --baud  
  Baud Rate

### Optional
*      --data  
  Number of Data Bits. default **8**
*      --parity  
  Parity Mode. `none`, `even` or `odd`. default **none**
*      --stop  
  Number of Stop Bits. default **1**

* -y / --binary  
  Binary Mode  
  parse input string as byte array

  ##### example
  ```
  0x56 0x78 0x9ABC

  is interpreted as  

  [0x56, 0x78, 0xBC, 0x9A]
  ```

### Others
* -l / --list  
  List COM Ports (**Windows Only**)
