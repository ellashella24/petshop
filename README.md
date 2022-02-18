<div id="top"></div>

# Petshop

<!-- PROJECT LOGO -->
<div align="center">
  <a href="https://github.com/ellashella24/petshop">
    <img src="images/logo.png" alt="Logo" width="180" height="180">
  </a>

  <h3 align="center">Petshop</h3>


  <p align="center">
    An E-Commerce App for Pet Shop
    <br />
    <div id = "other-software-design"></div>
    <a href="https://github.com/ellashella24/petshop/blob/main/documentation/wireflow.png?raw=true"">Wireflow</a>
    ·
    <a href="https://github.com/ellashella24/petshop/blob/main/documentation/usecase.jpeg">Use Case</a>
    ·
    <a href="https://github.com/ellashella24/petshop/blob/main/documentation/flowchart.jpeg">Flowchart</a>
    ·
    <a href="https://github.com/ellashella24/petshop/blob/main/documentation/erd.jpeg">ERD</a>
    ·
    <a href="https://app.swaggerhub.com/apis-docs/ellashella24/petshop/1.0.0">Open API</a>
  </p>
</div>
<br />

<!-- TABLE OF CONTENTS -->
## Table of Contents
1. [About the Project](#about-the-project)
2. [High Level Architecture](#high-level-architecture)
3. [Tech Stack](#tech-stack)
4. [Code Structure](#code-structure)
    - [Structuring](#structuring)
    - [Unit Test](#unit-test)
5. [How to Contrib](contribute.md)
6. [Contact](#contact)

<!-- ABOUT THE PROJECT -->
## About The Project
- An app that allow user to be a pet shop owner to sell their services and products or to be a customer to buy them. 
- Pet shop owners will be helped to market the products and services so that they can be easily reached by customers and they will be helped to get the products and services that are needed by their pets easily.
- Build with Golang, Echo Framework, MySQL adn GORM for manage repository, Xuri Excelize for Export List Product Selling to Excel, FTP to store Image Product to server, Xendit API for Payment Gateway, Deploy the project on [Okteto](https://ellashella24.cloud.okteto.net).

<p align="right">(<a href="#top">back to top</a>)</p>

## High Level Architecture

HLA design for this project shown in the picture below

<img src="images/HLA-rev-3.jpeg" alt="hla" width="800" height="462" >

<br />

<p align="right">(<a href="#top">back to top</a>)</p>

## Tech Stack
### RESTful-API
- [Go](https://go.dev/)
- [Echo Framework](https://echo.labstack.com/) - Go Framework
- [MySQL](https://www.mysql.com/) - SQL Database
- [GORM](https://gorm.io/index.html) - ORM Library
- [FTP](https://github.com/jlaffaye/ftp) - Upload File
- [SMTP](https://github.com/xhit/go-simple-mail) - Send Email
- [Xuri Excelize](https://xuri.me/excelize/) - Export Data to Excel Files
- [Xendit](https://www.xendit.co/id/?utm_source=google&utm_medium=cpc&utm_campaign=BKWS-Exact-ID-ID&utm_content=payment-gateway&utm_term=xendit) - Payment Gateway

### Deployment
- [Docker](https://www.docker.com/) - Container Images
- [Okteto](https://www.okteto.com/) - Kubernetes Platform
- [Kubernetes](https://kubernetes.io/) - Container 

Follow the link to see deployment flow of this project : [Deployment Flow](https://github.com/ellashella24/petshop/blob/main/documentation/documentation/deployment-flow.jpeg)

### Collaboration 
- [Trello](https://trello.com/) - Manage Project
- [Github](https://github.com/) - Versioning Project

<p align="right">(<a href="#top">back to top</a>)</p>

## Code Structure
This project use Layered Architure to organized each components into spesific function  

### Structuring
  ```sh
    petshop
    ├── config                        
    │     └──config.go                # Contains list of configuration of the project
    ├── constants                     
    │     └──constants.go             # Contains list constant variable
    ├── delivery                      # Contains list of component for handle request dan response
    │     └──common                   # Contains list of http request format based on the result from controller 
    │     │   ├── common.go           # Contains list of http request format
    │     │   └── http_responses.go   
    │     └──controller               # Contains list of component that receive the request and return a response
    │     │   ├── user
    │     │   ├── formatter_req.go    # Contains list of request format for each function on the controller
    │     │   ├── formatter_res.go    # Contains list of response format for each function on the controller
    │     │   ├── user_test.go        # Contains list of function for test each function on the controller
    │     │   └── users.go            # Contains list of controller for each entity    
    │     └──routes  
    │         └── routes.go           # Contains list of route to access each function on controller  
    ├── entity                        # Contains model all entity
    │     └── user.go                 # Contains model for spesific entity
    ├── repository                    # Contains list of functions that process the request and stores it in database
    │     ├── user_test.go            # Contains list of function for test each function on the repository
    │     └── users.go                # Contains list of repository for each entity
    ├── service                       # Contains list of function to access other components outside the project
    │     └── ftp.go                  # Containts list of function to manage files upload to the ftp
    │     └── export-excel.go         # Containts list of function to stores response data into the excel files
    ├── utils                         # Contains list of function to config each type of database
    │     └── mysqldriver.go          # Contains list of function to config MySQL type database
    ├── .env                          # Contains list of environment variable to run the project 
    ├── .gitignore                    # Contains list of directory/file name that will igonored when push project
    ├── go.mod                  
    ├── go.sum                  
    ├── main.go                       # Contains list of component that need to be executed first to run the app
    └── README.md    
  ```

### Unit Test
Coverage result on all functions is 99.2% which the most functions have reached 100% coverage. Coverage result for each function shown in the picture below

<img src="images/coverage-result-ver-2.jpg" alt="coverage-result">

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->
## Contact
* Naufal Aammar Hibatullah - [Github](https://github.com/nflhibatullah) · [LinkedIn](https://www.linkedin.com/in/naufal-hibatullah-441a58222/)
* Niendhitta Tamia Lassela - [Github](https://github.com/ellashella24) · [LinkedIn](https://www.linkedin.com/in/ntlassela/)

<p align="right">(<a href="#top">back to top</a>)</p>