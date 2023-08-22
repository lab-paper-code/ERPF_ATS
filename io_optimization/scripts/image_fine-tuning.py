import sys
import re

import tensorflow as tf
from tensorflow.keras.applications.mobilenet import MobileNet
from tensorflow.keras.layers import Dense, GlobalAveragePooling2D
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import Adam
import tensorflow.keras.preprocessing as preprocessing

from keras.callbacks import ModelCheckpoint

import os
# os.environ["CUDA_VISIBLE_DEVICES"] = "-1"

def get_args(argv: list) -> dict:
    arg_parser = re.compile(r'--(?P<option_name>[a-zA-Z_]+)=(?P<value>[\S]+)')
    valid_options = [
        'dataset_path',
        'dataset_size',
        'checkpoint_path',
        'log_path',
        'batch_size',
        'epochs',
        'base_model'
    ]
    settings = {
        # Default setting
        'dataset_path': './images',
        'dataset_size': 50000,
        'batch_size': 32,
        'epochs': 1,
        'base_model': None
    }

    for arg in argv:
        parsing_result = arg_parser.search(arg)
        if parsing_result is None:
            print('[ERR] Wrong option:', arg, '(Ignored)')
        else:
            option_name = parsing_result.group('option_name')
            option_value = parsing_result.group('value')
            if option_name not in valid_options:
                print('[ERR] Invalid option:', option_name, '(Ignored)')
            else:
                if option_name in ['dataset_size', 'batch_size', 'epochs']:
                    option_value = int(option_value)
                settings[option_name] = option_value
    
    return settings

def set_callbacks(settings:dict) -> list:
    callbacks = []
    
    if 'checkpoint_path' in settings.keys():
        if not settings['checkpoint_path'].endswith('/'): settings['checkpoint_path'] += '/'
        checkpoint_filename = 'checkpoint-epoch-{}.h5'.format(settings['epochs'])
        checkpoint = ModelCheckpoint(
            checkpoint_filename,
            monitor='accuracy',
            verbose=1,
            save_best_only=True,
            mode='auto'
            )
        callbacks.append(checkpoint)
    
    if 'log_path' in settings.keys():
        total_steps = settings['dataset_size']//settings['batch_size']
        # 중간의 100 Steps를 프로파일링
        mid = total_steps // 2
        start_step = mid - 49
        end_step = mid + 50
        if start_step <= 0:
            start_step = 1
        if end_step > total_steps:
            end_step = total_steps

        profile = tf.keras.callbacks.TensorBoard(
            log_dir=settings['log_path'],
            histogram_freq=1,
            profile_batch=[start_step, end_step]
        )
        callbacks.append(profile)
    
    return callbacks

def main():
    settings = get_args(sys.argv[1:])
    callbacks = set_callbacks(settings)

    # Set model
    base_model = MobileNet(weights=settings['base_model'], include_top=False, input_shape=(224, 224, 3))
    base_model.trainable = False
    inputs = base_model.input
    x = base_model.output
    x = GlobalAveragePooling2D()(x)
    predictions = Dense(1000, activation='softmax')(x)
    model = Model(inputs=inputs, outputs=predictions)
    model.compile(optimizer=Adam(learning_rate=0.0001), loss='categorical_crossentropy', metrics=['accuracy'])

    # Input pipeline:
    # Fetch → Preprocessing → Batch → Prefetch → Training
    '''
    data_gen = preprocessing.image.ImageDataGenerator(rescale=1./255)

    # Fetch + Preprocessing + Batch + Prefetch
    dataset = data_gen.flow_from_directory(
        settings['dataset_path'],
        batch_size=settings['batch_size'],
        class_mode='categorical'
    )
    '''

    # Fetch + Batch + Prefetch
    dataset = tf.keras.utils.image_dataset_from_directory(
      settings['dataset_path'],
      batch_size=settings['batch_size'],
      label_mode='categorical',
      image_size=(224, 224)
    )
    normalization_layer = tf.keras.layers.Rescaling(1./127.5, offset=-1)
    dataset = dataset.map(lambda x, y: (normalization_layer(x), y))

    # Training
    model.fit(
        dataset,
        epochs=settings['epochs'],
        callbacks=callbacks
    )

main()