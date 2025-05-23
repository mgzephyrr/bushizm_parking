import {type ChangeEvent, useEffect, useState} from 'react';
import TextField from '@mui/material/TextField';

interface IPhoneInputProps {
    phone: string
    setPhone: (value: string) => void
}

const PhoneInput = ({phone, setPhone}: IPhoneInputProps) => {
    const [cursorPosition, setCursorPosition] = useState(2);

    const handlePhoneChange = (e: ChangeEvent<HTMLInputElement>) => {
        const input = e.target.value;

        let cleaned = input.replace(/\D/g, '');

        if (cleaned.length > 11) {
            cleaned = cleaned.substring(0, 11);
        }

        if (cleaned.startsWith('7') && input.length === 1) {
            cleaned = cleaned.substring(1);
        }

        let formatted = '+7';
        if (cleaned.length > 1) {
            const rest = cleaned.substring(1);
            formatted += rest;
        }

        setPhone(formatted);

        const cursorPos = e.target.selectionStart || 2;
        setCursorPosition(Math.max(2, cursorPos));
    };

    useEffect(() => {
        const input = document.getElementById('phone-input') as HTMLInputElement;
        if (input) {
            input.setSelectionRange(cursorPosition, cursorPosition);
        }
    }, [phone, cursorPosition]);

    return (
        <TextField
            id="phone-input"
            label="Номер телефона"
            variant="outlined"
            value={phone}
            onChange={handlePhoneChange}
            inputProps={{
                maxLength: 16,
            }}
            sx={{width: "350px"}}
        />
    );
};

export default PhoneInput;