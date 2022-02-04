# Competitive programming site


# client_side
1️⃣ Each folder in this directory indicates problem id and contains solutions for the corresponding problem.

# cmd
2️⃣ Contains: 
 - main.go - entry point of the program
 - 2 go and exe files (participant_solution.go / participant_solution.exe and main_solution.go / main_solution.exe)
 - makefile that builds binaries for the participant_solution.go and main_solution.go

# internal
3️⃣ Server and the logic of handling requests.

# temp_solutions
4️⃣ Storage of users' solution files.

# web
5️⃣ The frontend part (templates and css). 

To run project on windows enter following commands from the root of the project:
```bash
cd deployments/kafka
docker compose up
```
Again on new terminal, from the root of the project:
```bash
cd cmd/myapp
go run main.go
```
To run frontend side:
```bash
cd web/frontend
go run main.go
```
To run makefile:
```bash
cd makeme
make run prog=path_to_prog exec=exec_name.exe
```
[Report](https://github.com/Kambar-ZH/Golang_Midterm_Project/blob/master/Report.pdf)