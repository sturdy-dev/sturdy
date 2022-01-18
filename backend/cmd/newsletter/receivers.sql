SELECT u.email
FROM users u
         LEFT JOIN notification_settings ns ON ns.user_id = u.id
WHERE u.email NOT LIKE '%@fake.getsturdy.com'
  AND u.email_verified = true
  AND u.email NOT LIKE '%foo%'
  AND u.email NOT LIKE '%bar%'
  AND u.email NOT LIKE '%+deleted@%'
  AND u.email NOT LIKE '%@gmail.clm'
  AND u.email != 'a@a.com'
  AND (ns.receive_newsletter = true OR ns.receive_newsletter IS NULL);