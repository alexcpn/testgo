# testgo
## Gocql Cassandra - Illustrating secondary indexes for pagination

# Short Summary

First there is a table with two fields. Just for illustration we use only few fields.

1.Say we insert a million rows with this

Along comes the product owner with a (rather strange) requirement that we need to list all the data as pages in the GUI. Assuming that there are hundred entries 10 pages each.

For this we update the table with a column called page_no.
Create a secondary index for this column.
Then do a one time update for this column with page numbers. Page number 10 will mean 10 contiguous rows updated with page_no as value 10.
Since we can query on a secondary index each page can be queried independently.
Code is self explanatory and here - https://github.com/alexcpn/testgo

Note caution on how to use secondary index properly abound. Please check it. In this use case I am hoping that i am using it properly. Have not tested with multiple clusters


# Results

--------------------------
Create and INSERT 1 million records

root@e6e1be5561d3:/coding/testgo# go install hello && ./bin/hello -ni 1000000 -nq 1000 -del=false
Hello, World -Test  Cassandra Pagination
INFO[0000] Number of Rows to Insert 1000000             
INFO[0000] Number of Rows to Query 1000                 
INFO[0000] Delete Table after tests false               
INFO[0000] Successfully connected                       
INFO[0682] insertRows Time took 11m22.341069466s        
INFO[0682] fetchRows Time took 61.246949ms 

-------------
Time to read 1 million 12 seconds - 

root@e6e1be5561d3:/coding/testgo# go install hello && ./bin/hello -ni 10 -nq 10000000 -del=false
Hello, World -Test  Cassandra Pagination
INFO[0000] Number of Rows to Insert 10                  
INFO[0000] Number of Rows to Query 10000000             
INFO[0000] Delete Table after tests false               
INFO[0000] Successfully connected                       
INFO[0000] insertRows Time took 9.56124ms               
INFO[0012] Number of rows read  1000000                 
INFO[0012] fetchRows Time took 12.277188868s            
Connected closed

--------------------

Created a column added a secondary index and updated page number

update took 12 mts

root@e6e1be5561d3:/coding/testgo# go install hello && ./bin/hello -ni 10 -nq 1000000 -del=false -pno=10
Hello, World -Test  Cassandra Pagination
INFO[0000] Number of Rows to Insert 10                  
INFO[0000] Number of Rows to Query 1000000              
INFrO[0000] Delete Table after tests false               
INFO[0000] Page Number to read 10                       
INFO[0000] Successfully connected                       
INFO[0000] insertRows Time took 9.120711ms              
INFO[0012] Number of rows read  1000000                 
INFO[0012] fetchRows Time took 12.140751116s            
WARN[0012] Could not alterTableStatement Table : Invalid column name page_no because it conflicts with an existing column 
INFO[0012] Altered Table                                
WARN[0012] Could not alterTableStatement Table : Index pageindx already exists 
INFO[0012] Created Index                                
INFO[0707] Number of rows read  1000000                 
INFO[0707] updateRows Time took 11m35.29930421s

.....
INFO[0707] Iter imsi: 133196636 7777133196636 adaddddddddadaasdass133196636   Page=10 
INFO[0707] Iter imsi: 133306343 7777133306343 adaddddddddadaasdass133306343   Page=10 
INFO[0707] Iter imsi: 132317741 7777132317741 adaddddddddadaasdass132317741   Page=10 
INFO[0707] Iter imsi: 132438485 7777132438485 adaddddddddadaasdass132438485   Page=10 
INFO[0707] Iter imsi: 132852300 7777132852300 adaddddddddadaasdass132852300   Page=10 
INFO[0707] Iter imsi: 132626942 7777132626942 adaddddddddadaasdass132626942   Page=10 
INFO[0707] Iter imsi: 132631436 7777132631436 adaddddddddadaasdass132631436   Page=10 
INFO[0707] Iter imsi: 132649592 7777132649592 adaddddddddadaasdass132649592   Page=10 
INFO[0707] Iter imsi: 132814138 7777132814138 adaddddddddadaasdass132814138   Page=10 
INFO[0707] Iter imsi: 133084740 7777133084740 adaddddddddadaasdass133084740   Page=10 
INFO[0707] Iter imsi: 132510301 7777132510301 adaddddddddadaasdass132510301   Page=10 
INFO[0707] Iter imsi: 132903438 7777132903438 adaddddddddadaasdass132903438   Page=10 
INFO[0707] Iter imsi: 133073319 7777133073319 adaddddddddadaasdass133073319   Page=10 
INFO[0707] Iter imsi: 133113857 7777133113857 adaddddddddadaasdass133113857   Page=10 
INFO[0707] Iter imsi: 132479935 7777132479935 adaddddddddadaasdass132479935   Page=10 
INFO[0707] Iter imsi: 132443965 7777132443965 adaddddddddadaasdass132443965   Page=10 
INFO[0707] Number of rows read  0                       
INFO[0707] fetchRows By Page Time took 440.224143ms     
Connected closed
root@e6e1be5561d3:/coding/testgo# 

-------------
Query by the secondary index now

/coding/testgo# go install hello && ./bin/hello -ni 10 -nq 10 -del=false -pno=20
INFO[0000] Iter imsi: 133099417 7777133099417 adaddddddddadaasdass133099417   Page=20 
INFO[0000] Iter imsi: 132766587 7777132766587 adaddddddddadaasdass132766587   Page=20 
INFO[0000] Iter imsi: 132803921 7777132803921 adaddddddddadaasdass132803921   Page=20 
INFO[0000] Iter imsi: 132537300 7777132537300 adaddddddddadaasdass132537300   Page=20 
INFO[0000] Number of rows read  0                       
INFO[0000] fetchRows By Page Time took 153.256121ms   



