# Run docker with cmd: docker run -d -p 8888:8080 [image id]

FROM ubuntu:latest

RUN mkdir Module3
COPY API-Gateway/main /Module3/App/API-Gateway/
COPY Core/Merchant/Manage_Account/main /Module3/App/Core/Merchant/Manage_Account/
COPY Core/Merchant/Manage_Bill/main /Module3/App/Core/Merchant/Manage_Bill/
COPY Core/Merchant/Manage_Support/main /Module3/App/Core/Merchant/Manage_Support/
COPY Middle-ware/main /Module3/App/Middle-ware/


CMD ["./Module3/App/Core/Merchant/Manage_Support/main","./Module3/App/Core/Merchant/Manage_Bill/main","./Module3/App/Core/Merchant/Manage_Account/main","./Module3/App/API-Gateway/main","./Module3/App/Middle-ware/main"]
# RUN ./Module3/App/Core/Merchant/Manage_Bill/main -D
# RUN ./Module3/App/Core/Merchant/Manage_Account/main -D
# RUN ./Module3/App/API-Gateway/main -D
# RUN ./Module3/App/Middle-ware/main -D
