# GOKIT Project Structure

GOKIT project aims to support fast-driven application development using the creational pattern with the factory method. GOKIT provides a structured way to build and manage services, making it easy to integrate with other services, and it is implemented with a standardized OpenTelemetry observer.

## Feature Builded
[v] **Adapter** Already supported databases (sql/db, gorm, redis)

[v] **Logger** Already built-in log for all kind Protocol (gRPC, REST) and outgoing request log

[v] **REST Protocol** Already supported rest 
protocol with Go-Fiber

[v] **Grpc Protocol** Already supported rest protocol with Go GRPC

[v] **Standarize Error** Standarize Error Handling with errorkit

## Future Feature
[ ] **Broker Support** Comming soon

[ ] **Tracer** Comming soon, supported Open Telemetry Collector


## Directories Structure

- **abstract**: This folder contains abstract definitions and interfaces that can be implemented by different components of the project.

- **adapter**: This folder is used for the adapter pattern, providing an interface to interact with various services or libraries, facilitating communication between different parts of the application.

- **config**: This folder holds configuration files and settings required to run the application. It may include environment variable management and configuration loading logic.

- **factory**: This folder contains factory methods that create instances of various types. It abstracts the instantiation process, allowing for more flexible and maintainable code.

- **logger**: This folder includes logging utilities to capture and manage log messages throughout the application. It helps in monitoring and debugging.

- **protoc**: This folder is likely dedicated to Protocol Buffers definitions and configurations, enabling the application to define and serialize structured data.

- **tracer**: This folder contains code related to tracing functionality, which is used for monitoring and analyzing the execution flow of the application.

- **types**: This folder defines various data types used throughout the application, including structures and enumerations.

- **utils**: This folder contains utility functions and helpers that provide common functionalities used across different parts of the application.