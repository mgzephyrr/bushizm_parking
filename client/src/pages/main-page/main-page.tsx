import Header from "../../components/header/header.tsx";
import Footer from "../../components/footer/footer.tsx";

export default function MainPage() {
    return (
        <>
            <Header/>
            <div style={{height: "100%"}}>
                <div style={{display: "flex", justifyContent: "center", alignItems: "center"}}>
                    <button>Занять место</button>
                </div>

            </div>
            <Footer>
                <span>Большой бушизм ©️ 2025</span>
            </Footer>
        </>
    )
};