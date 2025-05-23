import {Box, FormControl, Modal} from "@mui/material";
import PhoneInput from "../phone-input/phone-input.tsx";
import {useState} from "react";
import style from "./modal-form.module.css"
import {useDispatch} from "react-redux";
import {login} from "../../app/api/slice/user-slice.ts";
import type {ILoginResponse} from "../../app/api/auth/types.ts";

interface IModalFormProps {
    isOpen: boolean;
    handleClose: () => void;
    trigger: any
}

export default function ModalForm({isOpen, handleClose, trigger}: IModalFormProps) {
    const dispatch = useDispatch()
    const [phone, setPhone] = useState<string>("+7")

    const handleFormSubmit = () => {
        trigger({phone: phone.substring(1, 12)}).then((result: { data: ILoginResponse }) => {
            dispatch(login(result?.data?.full_name))
        }).then(() => handleClose())
    }
    return (
        <Modal open={isOpen} onClose={handleClose}>
            <Box className={style.form} sx={style}>
                <h2 style={{color: "black"}}>Авторизация</h2>
                <FormControl sx={{display: "flex", flexDirection: "column", gap: "12px", padding: "24px"}}>
                    <PhoneInput phone={phone} setPhone={setPhone}/>
                    <button className={style.formButton} onClick={handleFormSubmit}>Отправить</button>
                </FormControl>
            </Box>
        </Modal>
    )
}