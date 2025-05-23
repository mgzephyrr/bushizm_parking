import {useGetSpotsNumberQuery, useLazyGetSubscribeQuery} from "../../app/api/subscriptions/subscriptions.ts";
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
import {LogInIcon} from "../../assets/icons/LogInIcon.tsx";
import {LogOutIcon} from "../../assets/icons/LogoutIcon.tsx";
import {UjinLogo} from "../../assets/icons/UjinLogo.tsx";

export default function MainPage() {
    const userName = useSelector((state: RootState) => state.user.full_name)

    const dispatch = useDispatch();

    const {data: spotsNumberData, refetch: refetchSpots} = useGetSpotsNumberQuery(null)
    const {data, refetch: refetchNotify} = useGetNotifyQuery(null)
    const [isAlert, setIsAlert] = useState(false)
    const [trigger] = useLazyGetSubscribeQuery()
    const [loginTrigger] = useGetLoginMutation()
    const [logoutTrigger] = useGetLogoutMutation()

    const [isModalOpen, setIsModalOpen] = useState(false);

      useEffect(() => {
        const interval = setInterval(async () => {
          try {
            const { data } = await refetchNotify(); // Выполняем запрос
            refetchSpots()
            if (data?.send === "yes") {
              alert("Для вас освободилось место на парковке");
            }
          } catch (error) {
            console.error("Ошибка при проверке уведомлений:", error);
          }
        }, 3000);

        return () => clearInterval(interval);
      }, [refetchNotify]); // Зависимость только от refetchNotify

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
                    <UjinLogo/>
                    <div className={styles.headerUserName}>{userName ? userName : ""}</div>
                    <button onClick={userName ? handleLogout : handleModalOpen} className={styles.loginButton}>
                        {!userName ? <LogInIcon/> : <LogOutIcon/>}
                    </button>
                </div>
            }/>
            <div className={styles.mainSectionWrapper}>
                <div className={styles.countSection}>
                    <div>{spotsNumberData?.spots_number ?? 0} свободных мест</div>
                </div>

                <div className={styles.buttonSection}>
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
