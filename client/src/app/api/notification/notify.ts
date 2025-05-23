import {fetchBaseQuery} from "@reduxjs/toolkit/query";
import {createApi} from "@reduxjs/toolkit/query/react";


export const notificationApi = createApi({
    reducerPath: 'notificationApi',
    baseQuery: fetchBaseQuery({baseUrl: 'http://127.0.0.1:8001/', credentials: "include",}),
    endpoints: (build) => ({
        getNotify: build.query<{ send: string }, string>({
            query: (accessToken) => ({
                url: '/notify',
                body: {token: accessToken},
                method: 'POST',
            })
        }),
    }),
})

export const {useGetNotifyQuery} = notificationApi