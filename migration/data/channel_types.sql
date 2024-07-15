insert into channel_type(id, k, name) VALUES
        ('01', 'Mail', 'Email Notification channel') on conflict(id) do nothing;

insert into channel_type(id, k, name) VALUES
    ('02', 'SMS', 'GMS SMS Notification') on conflict(id) do nothing;

insert into channel_type(id, k, name) VALUES
    ('03', 'GCM', 'Google Cloud Messaging') on conflict(id) do nothing;