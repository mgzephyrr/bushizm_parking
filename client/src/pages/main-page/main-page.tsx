import {useLazyGetSubscribeQuery} from "../../app/api/subscriptions/subscriptions.ts";
import ModalForm from "../../components/modal-form/modal-form.tsx";
import {useEffect, useState} from "react";
import styles from "./main-page.module.css"
import Header from "../../components/header/header.tsx";
import Footer from "../../components/footer/footer.tsx";
import {useGetLoginMutation, useGetLogoutMutation} from "../../app/api/auth/auth.ts";
import {useDispatch, useSelector} from "react-redux";
import type {RootState} from "../../app/api/store.ts";
import {logout} from "../../app/api/slice/user-slice.ts";
import {useGetNotifyQuery} from "../../app/api/notification/notify.ts";

export default function MainPage() {
    const userName = useSelector((state: RootState) => state.user.full_name)
    const [notifyState, setNotifyState] = useState("")

    const dispatch = useDispatch();
    const {data, refetch: refetchNotify} = useGetNotifyQuery(null)
    const [trigger] = useLazyGetSubscribeQuery()
    const [loginTrigger] = useGetLoginMutation()
    const [logoutTrigger] = useGetLogoutMutation()

    const [isModalOpen, setIsModalOpen] = useState(false);

    useEffect(() => {
        const interval = setInterval(() => {
            refetchNotify()
            if (!notifyState) {
                setNotifyState(data?.send ?? "")
            }
            if (notifyState && notifyState !== data?.send) {
                alert("Можно парковаца)))")
            }
        }, 15000);

        return () => clearInterval(interval);
    }, [refetchNotify]);

    const handleModalClose = () => {
        setIsModalOpen(false);
    }

    const handleModalOpen = () => {
        setIsModalOpen(true);
    }

    const handleButtonClick = () => {
        trigger(null)
    }

    const handleLogout = () => {
        logoutTrigger({})
        dispatch(logout())
    }

    return (
        <>
            <Header children={
                <div className={styles.headerContent}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="144" height="51" viewBox="0 0 144 51" fill="none">
                        <path
                            d="M88.0866 10.707H81.381V33.804L75.4111 33.804C71.4256 33.804 68.7149 31.0879 68.7149 27.1V15.8648C68.7149 13.0164 66.4071 10.7072 63.5603 10.7072H62.0093L62.0093 27.1C62.0093 34.506 68.0095 40.5098 75.4111 40.5098L88.0866 40.5096V10.707Z"
                            fill="#212121"/>
                        <path
                            d="M117.146 40.5103L123.852 40.5103L123.852 17.4132L129.822 17.4133C133.807 17.4133 136.518 20.1294 136.518 24.1172L136.518 35.3525C136.518 38.2009 138.826 40.5101 141.673 40.5101L143.224 40.5101L143.224 24.1172C143.224 16.7112 137.223 10.7075 129.822 10.7075L117.146 10.7077L117.146 40.5103Z"
                            fill="#212121"/>
                        <path
                            d="M111.964 7.375C111.964 6.49441 111.79 5.62244 111.453 4.80888C111.116 3.99532 110.622 3.25611 110 2.63343C109.377 2.01076 108.638 1.51683 107.824 1.17985C107.011 0.84286 106.139 0.669415 105.258 0.669415L105.258 7.375L111.964 7.375Z"
                            fill="#212121"/>
                        <path
                            d="M100.016 43.8274C100.016 44.708 99.8428 45.5799 99.5058 46.3935C99.1688 47.2071 98.6749 47.9463 98.0522 48.569C97.4296 49.1916 96.6904 49.6856 95.8768 50.0225C95.0632 50.3595 94.1913 50.533 93.3107 50.533L93.259 10.6885L99.9646 10.6885L100.016 43.8274Z"
                            fill="#212121"/>
                        <path d="M105.216 10.707H111.922V40.5096H105.216V10.707Z" fill="#212121"/>
                        <path
                            d="M0 7.72125C0 3.81127 3.16967 0.641602 7.07965 0.641602H50V50.6416H7.07965C3.16967 50.6416 0 47.4719 0 43.562V7.72125Z"
                            fill="#00C0C9"/>
                        <path
                            d="M41.298 7.42627H32.9162V35.447L25.454 35.4469C20.4723 35.4469 17.0839 32.0405 17.0839 27.0392V13.8949C17.0839 10.3225 14.1993 7.42649 10.6409 7.42649H8.70215V27.0392C8.70215 36.3274 16.2023 43.8569 25.454 43.8569L41.298 43.8567V7.42627Z"
                            fill="white"/>
                    </svg>
                    <div className={styles.headerUserName}>{userName ? userName : ""}</div>
                    <button onClick={userName ? handleLogout : handleModalOpen} className={styles.loginButton}>
                        {!userName ? (
                                <svg xmlns="http://www.w3.org/2000/svg" width="35" height="35" viewBox="0 0 24 24">
                                    <path d="m13 16l5-4l-5-4v3H4v2h9z"/>
                                    <path
                                        d="M20 3h-9c-1.103 0-2 .897-2 2v4h2V5h9v14h-9v-4H9v4c0 1.103.897 2 2 2h9c1.103 0 2-.897 2-2V5c0-1.103-.897-2-2-2"/>
                                </svg>)
                            : (<svg xmlns="http://www.w3.org/2000/svg" width="35" height="35" viewBox="0 0 512 512">
                                <path
                                    d="M77.155 272.034H351.75v-32.001H77.155l75.053-75.053v-.001l-22.628-22.626l-113.681 113.68l.001.001h-.001L129.58 369.715l22.628-22.627v-.001z"/>
                                <path d="M160 16v32h304v416H160v32h336V16z"/>
                            </svg>)}
                    </button>
                </div>
            }/>
            <div className={styles.mainSectionWrapper}>
                <div className={styles.countSection}>
                    <div>12 свободных мест</div>
                    <div>0 людей в очереди</div>
                </div>

                <div style={{display: "flex", justifyContent: "center", alignItems: "center"}}>
                    <button className={styles.button} onClick={handleButtonClick}>
                        Занять место
                    </button>
                </div>
            </div>
            <Footer>
                <span>БОЛЬШОЙ БУШИЗМ ©️ 2025</span>
            </Footer>
            <ModalForm trigger={loginTrigger} isOpen={isModalOpen} handleClose={handleModalClose}/>
        </>
    )
};