USE music_transfer;

DELETE FROM app.playlists WHERE name LIKE 'TEST_%';
DELETE FROM app.users WHERE email LIKE 'TEST_%';