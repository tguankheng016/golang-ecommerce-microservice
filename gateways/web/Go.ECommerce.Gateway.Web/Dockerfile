FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
USER $APP_UID
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["Go.ECommerce.Gateway.Web/Go.ECommerce.Gateway.Web.csproj", "Go.ECommerce.Gateway.Web/"]
RUN dotnet restore "Go.ECommerce.Gateway.Web/Go.ECommerce.Gateway.Web.csproj"
COPY . .
WORKDIR "/src/Go.ECommerce.Gateway.Web"
RUN dotnet build "Go.ECommerce.Gateway.Web.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "Go.ECommerce.Gateway.Web.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Go.ECommerce.Gateway.Web.dll"]

