import { ForgotPasswordForm } from '@/components';
import styles from '../login/login.module.css';
import Image from 'next/image';

export default function ForgotUsernamePage() {
  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.logo}>
          <Image src="/icons/logo.svg" alt="Logo" width={180} height={100} />
        </div>

        <ForgotPasswordForm type="username" />
      </div>
    </div>
  );
}
