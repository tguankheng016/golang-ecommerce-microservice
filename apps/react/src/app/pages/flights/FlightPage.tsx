import DefaultPage, { BreadcrumbItem } from "@app/components/layout/DefaultPage";
import APIClient from "@shared/service-proxies/api-client";
import { useEffect } from "react";

const FlightPage = () => {
    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Flights')
    ];

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        const flightService = APIClient.getFlightService();

        flightService.getFlights(undefined, undefined, undefined, undefined, undefined, 0, 10, '', '', signal)
            .then((res) => {
                console.log(res);
            });

        return () => {
            abortController.abort();
        };
    }, []);
    
    return (
        <DefaultPage title="Flights" breadcrumbs={breadcrumbs}>
            <div className="card">
                <div className="card-body p-lg-17">
                    <p>
                        FlightPage
                    </p>
                </div>
            </div>
        </DefaultPage>
    )
}

export default FlightPage