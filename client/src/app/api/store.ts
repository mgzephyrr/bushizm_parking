import {configureStore} from "@reduxjs/toolkit";
import {subscriptionsApi} from "./subscriptions/subscriptions.ts";
import {authApi} from "./auth/auth.ts";
import userReducer from "../api/slice/user-slice.ts"
import {notificationApi} from "./notification/notify.ts";

export const store = configureStore({
    reducer: {
        [subscriptionsApi.reducerPath]: subscriptionsApi.reducer,
        [authApi.reducerPath]: authApi.reducer,
        [notificationApi.reducerPath]: notificationApi.reducer,
        user: userReducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware()
            .concat(subscriptionsApi.middleware)
            .concat(authApi.middleware)
            .concat(notificationApi.middleware),

})

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;