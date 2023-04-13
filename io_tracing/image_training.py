import sys

import tensorflow as tf
from tensorflow.keras.applications.mobilenet import MobileNet
from tensorflow.keras.layers import Dense, GlobalAveragePooling2D
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import Adam
from tensorflow.keras.preprocessing.image import ImageDataGenerator

if len(sys.argv) != 4:
  print('Usage: python3 image_training.py <number_of_images> <batch_size> <epochs>')

n = int(sys.argv[1])
batch_size = int(sys.argv[2])
epochs = int(sys.argv[3])

# Load MobileNet model without top layers
base_model = MobileNet(weights='imagenet', include_top=False, input_shape=(224, 224, 3))

# Add top layers to the model
x = base_model.output
x = GlobalAveragePooling2D()(x)
x = Dense(1024, activation='relu')(x)
predictions = Dense(1000, activation='softmax')(x)
model = Model(inputs=base_model.input, outputs=predictions)

# Freeze the base model layers
for layer in base_model.layers:
    layer.trainable = False

# Compile the model
model.compile(optimizer=Adam(lr=0.0001), loss='categorical_crossentropy', metrics=['accuracy'])

# Load and preprocess the data
train_datagen = ImageDataGenerator(rescale=1./255)
train_generator = train_datagen.flow_from_directory(
        './images/validation_images',
        target_size=(224, 224),
        batch_size=batch_size,
        class_mode='categorical')

# Fine-tune the model
model.fit(
        train_generator,
        steps_per_epoch=n//batch_size, # number of images in the dataset divided by the batch size
        epochs=epochs)

# Save the model
model.save('mobile_net_finetuned.h5')
