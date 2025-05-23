import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import MainPage from "./pages/main-page/main-page.tsx";
import {Provider} from "react-redux";
import {store} from "./app/api/store.ts";

const router = createBrowserRouter([
    {
        path: "/",
        element: <MainPage/>,
    },
]);

createRoot(document.getElementById('root')!).render(
    <Provider store={store}>
        <StrictMode>
            <RouterProvider router={router}/>
        </StrictMode>
    </Provider>,
)
