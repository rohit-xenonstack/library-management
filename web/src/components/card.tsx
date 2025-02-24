import clsx from 'clsx'
import type { ReactNode } from 'react'
import type React from 'react'

import styles from '../styles/modules/card.module.scss'

interface CardProps {
  children?: ReactNode
  header?: ReactNode
  footer?: ReactNode
  className?: string
  variant?: 'primary' | 'secondary' // Optional variants
}

const Card: React.FC<CardProps> = ({
  children,
  header,
  footer,
  className,
  variant,
}) => {
  const cardClasses = clsx(
    styles.card,
    className,
    variant && styles[`card--${variant}`],
  )

  return (
    <div className={cardClasses}>
      {header && <div className={styles.cardHeader}>{header}</div>}
      <div className={styles.cardContent}>{children}</div>
      {footer && <div className={styles.cardFooter}>{footer}</div>}
    </div>
  )
}

export default Card
