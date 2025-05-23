import {fetchBaseQuery} from "@reduxjs/toolkit/query";
import {createApi} from "@reduxjs/toolkit/query/react";
import type {ILoginRequest, ILoginResponse} from "./types.ts";

export const authApi = createApi({
    reducerPath: 'authApi',
    baseQuery: fetchBaseQuery({baseUrl: 'http://127.0.0.1:8000/', credentials: "include",}),
    endpoints: (build) => ({
        getLogin: build.mutation<ILoginResponse, ILoginRequest>({
            query: (credentials) => ({
                url: '/login',
                method: 'POST',
                body: credentials
            })
        }),
        getLogout: build.mutation<{ message: string }, {}>({
            query: () => ({
                url: '/logout',
                method: 'POST',
            })
        }),
    }),
})

export const {useGetLoginMutation, useGetLogoutMutation} = authApi