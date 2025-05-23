import {fetchBaseQuery} from "@reduxjs/toolkit/query";
import {createApi} from "@reduxjs/toolkit/query/react";

export const subscriptionsApi = createApi({
    reducerPath: 'subscriptionsApi',
    baseQuery: fetchBaseQuery({baseUrl: 'http://localhost:8080/api/v1/', credentials: "include"}),
    endpoints: (build) => ({
        getSubscribe: build.query<null, null>({
            query: () => ({
                url: '/subscriptions/subscribe',
                method: 'POST',
            })
        }),
        getSpotsNumber: build.query<{ spots_number: number }, null>({
            query: () => ({
                url: '/spotsnumber',
            })
        }),
    }),
})

export const {useLazyGetSubscribeQuery, useGetSpotsNumberQuery} = subscriptionsApi
