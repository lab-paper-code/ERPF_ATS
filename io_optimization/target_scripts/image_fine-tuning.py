import sys

import tensorflow as tf
from tensorflow.keras.applications.mobilenet import MobileNet
from tensorflow.keras.layers import Dense, GlobalAveragePooling2D
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import Adam
import tensorflow.keras.preprocessing as preprocessing

if len(sys.argv) != 5:
  print('Usage: python3 image_fine-tuning.py <dataset_path> <dataset_size> <batch_size> <epochs>')

dataset_path = sys.argv[1]
dataset_size = int(sys.argv[2])
batch_size = int(sys.argv[3])
epochs = int(sys.argv[4])

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
model.compile(optimizer=Adam(learning_rate=0.0001), loss='categorical_crossentropy', metrics=['accuracy'])

# Load and preprocess the data
dataset = preprocessing.image_dataset_from_directory(
        dataset_path,
        image_size=(224, 224),
        batch_size=batch_size,
        labels='inferred',
        label_mode='categorical')

# Fine-tune the model
model.fit(
        dataset,
        steps_per_epoch=dataset_size//batch_size, # number of images in the dataset divided by the batch size
        epochs=epochs)