import React from "react";

interface FooterProps {
    children?: React.ReactNode;
}

export default function Footer({children}: FooterProps) {
    return <div style={{
        // position: "absolute",
        // bottom: 0,
        // left: 0,
        // right: 0,
        display: "flex",
        backgroundColor: "#569EB7",
        height: "70px",
        justifyContent: "center",
        alignItems: 'center'
    }}>{children}</div>;
}