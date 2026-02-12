extends CharacterBody2D

var socket = WebSocketPeer.new()
@onready var animated_sprite_2d = $AnimatedSprite2D

@onready var line_edit = $CanvasLayer/Control/LineEdit
@onready var chat_log = $CanvasLayer/Control/RichTextLabel

func _ready():
	socket.connect_to_url("ws://localhost:8080/ws")
	line_edit.visible = false

func _process(_delta):
	socket.poll()
	if socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		var input = {"key": "", "msg": ""}
  
		if not line_edit.has_focus():
			if Input.is_action_pressed("move_right"):
				input["key"] = "right"
				animated_sprite_2d.play("Walking")
				animated_sprite_2d.flip_h = false
			elif Input.is_action_pressed("move_left"):
				input["key"] = "left"
				animated_sprite_2d.play("Walking")
				animated_sprite_2d.flip_h = true
			elif Input.is_action_pressed("move_up"):
				input["key"] = "up"
				animated_sprite_2d.play("Walking")
			elif Input.is_action_pressed("move_down"):
				input["key"] = "down"
				animated_sprite_2d.play("Walking")
			else:
				animated_sprite_2d.play("IDLE")
  
		socket.send_text(JSON.stringify(input))
  
		while socket.get_available_packet_count() > 0:
			var packet = socket.get_packet()
			var data = JSON.parse_string(packet.get_string_from_utf8())
   
			if data:
				var target_pos = Vector2(data["X"], data["Y"])
			
				var direction = target_pos - position
				velocity = direction * 20
	
				move_and_slide()
	
			if data["Msg"] != "":
				chat_log.append_text("\n" + data["Msg"])

func _input(event):
	if event is InputEventKey and event.pressed and event.keycode == KEY_T:
		line_edit.visible = !line_edit.visible
		if line_edit.visible:
			line_edit.grab_focus()
		else:
			line_edit.release_focus()

	if event.is_action_pressed("ui_accept") and line_edit.has_focus():
		if line_edit.text != "":
			var chat_data = {"key": "", "msg": "Бро: " + line_edit.text}
			socket.send_text(JSON.stringify(chat_data))
			line_edit.text = ""
		line_edit.visible = false
		line_edit.release_focus()
