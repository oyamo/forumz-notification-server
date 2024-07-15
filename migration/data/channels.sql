insert into channel(id, channel_type_id, name)
    values ('0190ab09-e093-7e27-894c-945c1717e499', '01', 'SendGrid') on conflict (id) do nothing;

insert into channel(id, channel_type_id, name)
    values ('0190ab0a-31c7-7ae0-8b3a-a4b1f9936449', '01', 'Mailtrap') on conflict (id) do nothing;