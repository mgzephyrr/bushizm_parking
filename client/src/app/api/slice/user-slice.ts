import {createSlice} from "@reduxjs/toolkit";

const userSlice = createSlice({
    name: 'user',
    initialState: {
        full_name: ""
    },
    reducers: {
        login: (state, action) => {
            state.full_name = action.payload
        },
        logout: state => {
            state.full_name = ""
        }
    }
})

export const {login, logout} = userSlice.actions

export default userSlice.reducer
